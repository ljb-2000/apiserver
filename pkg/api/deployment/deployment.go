func DeploymentApp(req *http.Request) (code string, ret interface{}) {
	decoder := json.NewDecoder(req.Body)
	app := &dao.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("请求参数有误：%s", err.Error())
		code = StatusBadRequest
		ret = map[string]interface{}{"success": false, "reason": "请求参数有误，请检查！"}
		return
	}

	app.Status = dao.AppBuilding
	if err = app.Insert(); err != nil {
		log.Errorf("插入数据库失败：%s", err.Error())
		code = StatusInternalServerError
		ret = map[string]interface{}{"success": false, "reason": "插入数据库失败！"}
		return
	}

	//create a namespace
	nc := new(v1.Namespace)
	ncTypeMeta := unversioned.TypeMeta{Kind: "NameSpace", APIVersion: "v1"}
	nc.TypeMeta = ncTypeMeta

	nc.ObjectMeta = v1.ObjectMeta{
		Name: app.UserName,
	}

	nc.Spec = v1.NamespaceSpec{}

	_, err = dao.Clientset.Core().Namespaces().Create(nc)

	if err != nil {
		log.Errorf("deploy application failed ,the reason is %s", err.Error())
		app.Status = dao.AppFailed
		if err = app.Update(); err != nil {
			log.Errorf("update application status failed,the reason is %s", err.Error())
		}
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	//create a replicationController
	rc := new(v1.ReplicationController)

	rcTypeMeta := unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}
	rc.TypeMeta = rcTypeMeta

	rcObjectMeta := v1.ObjectMeta{Name: app.Name, Namespace: app.UserName, Labels: map[string]string{"name": app.Name}}
	rc.ObjectMeta = rcObjectMeta

	rcSpec := v1.ReplicationControllerSpec{
		Replicas: &app.InstanceCount,
		Selector: map[string]string{
			"name": app.Name,
		},
		Template: &v1.PodTemplateSpec{
			v1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.UserName,
				Labels: map[string]string{
					"name": app.Name,
				},
			},
			v1.PodSpec{
				RestartPolicy: v1.RestartPolicyAlways,
				Containers: []v1.Container{
					v1.Container{
						Name:  app.Name,
						Image: app.Image,
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								ContainerPort: 9080,
								Protocol:      v1.ProtocolTCP,
							},
						},
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse(app.Cpu),
								v1.ResourceMemory: resource.MustParse(app.Memory),
							},
						},
					},
				},
			},
		},
	}
	rc.Spec = rcSpec

	result, err := dao.Clientset.Core().ReplicationControllers(NameSpace).Create(rc)
	if err != nil {
		log.Errorf("deploy application failed ,the reason is %s", err.Error())
		app.Status = dao.AppFailed
		if err = app.Update(); err != nil {
			log.Errorf("update application status failed,the reason is %s", err.Error())
		}
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	} else {
		//create service
		service := new(v1.Service)

		svTypemeta := unversioned.TypeMeta{Kind: "Service", APIVersion: "v1"}
		service.TypeMeta = svTypemeta

		svObjectMeta := v1.ObjectMeta{Name: app.Name, Namespace: app.UserName, Labels: map[string]string{"name": app.Name}}
		service.ObjectMeta = svObjectMeta

		svServiceSpec := v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name:       app.Name,
					Port:       9080,
					TargetPort: intstr.FromInt(9080),
					Protocol:   "TCP",
					// NodePort:   32107,
				},
			},
			Selector: map[string]string{"name": app.Name},
			Type:     v1.ServiceTypeNodePort,
			// LoadBalancerIP: "172.17.11.2",
			// Status: v1.ServiceStatus{
			// 	LoadBalancer: v1.LoadBalancerStatus{
			// 		Ingress: []v1.LoadBalancerIngress{
			// 			v1.LoadBalancerIngress{IP: "172.17.11.2"},
			// 		},
			// 	},
			// },
		}
		service.Spec = svServiceSpec

		_, err := dao.Clientset.Core().Services(NameSpace).Create(service)

		if err != nil {
			log.Errorf("deploy application failed ,the reason is %s", err.Error())
			app.Status = dao.AppFailed
			if err = app.Update(); err != nil {
				log.Errorf("update application status failed,the reason is %s", err.Error())
			}
			code = StatusInternalServerError
			ret = JSON_EMPTY_OBJ
			return
		}
	}

	app.Status = dao.AppSuccessed
	if err = app.Update(); err != nil {
		log.Errorf("update application status failed,the reason is %s", err.Error())
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = result
	return
}

func Delete(req *http.Request) (code string, ret interface{}) {
	decoder := json.NewDecoder(req.Body)
	app := &dao.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("请求参数有误：%s", err.Error())
		code = StatusBadRequest
		ret = map[string]interface{}{"success": false, "reason": "请求参数有误，请检查！"}
		return
	}

	// rc := new(v1.ReplicationController)
	NameSpace = app.UserName
	rc, err := dao.Clientset.Core().ReplicationControllers(NameSpace).Get(app.Name)
	if err != nil {
		code = StatusNotFound
		ret = map[string]interface{}{"success": false, "reason": err.Error()}
		return
	}

	//delete conditions
	deleteOption := new(api.DeleteOptions)
	deleteOption.TypeMeta = unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}

	//if the value is 0 ,delete immediately. if not set,the default grace period for the specified type will be used ,
	//the default grace period is 30s
	deleteOption.GracePeriodSeconds = new(int64)

	//delete precondition(前提条件)
	deleteOption.Preconditions = &api.Preconditions{UID: &(rc.ObjectMeta.UID)}

	// If true/false，  added to/removed from the object's finalizers list
	deleteOption.OrphanDependents = parseUtil.BoolToPointer(false)
	err = dao.Clientset.Core().ReplicationControllers(NameSpace).Delete(app.Name, deleteOption)
	if err != nil {
		log.Errorf("delete application failed ：%s", err.Error())
		code = StatusInternalServerError
		ret = map[string]interface{}{"success": false, "reason": err.Error()}
		return
	}

	if err = app.Delete(); err != nil {
		log.Errorf("delete application failed,the reason is %s", err.Error())
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusNoContent
	ret = OK
	return
}

func Update(req *http.Request) (code string, ret interface{}) {
	vebrType := req.FormValue("vebrType")

	decoder := json.NewDecoder(req.Body)
	app := &dao.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("请求参数有误：%s", err.Error())
		code = StatusBadRequest
		ret = map[string]interface{}{"success": false, "reason": "请求参数有误，请检查！"}
		return
	}

	NameSpace = app.UserName
	rc := new(v1.ReplicationController)
	rc, _ = dao.Clientset.Core().ReplicationControllers(NameSpace).Get(app.Name)
	//stop application
	if vebrType == "stop" {
		rc.Spec.Replicas = parseUtil.Int32ToPointer(0)
		app.UpdateStatus = dao.StopFailed
	}

	//start application
	if vebrType == "start" {
		app, err := app.QueryOne() //cording  appname rcname to query app record
		if err != nil {
			log.Errorf("query application failed ：%s", err.Error())
			code = StatusInternalServerError
			ret = map[string]interface{}{"success": false, "reason": err.Error()}
			return
		}
		rc.Spec.Replicas = parseUtil.Int32ToPointer(app.InstanceCount)
		app.UpdateStatus = dao.StartFailed
	}

	//scale application
	if vebrType == "scale" {
		rc.Spec.Replicas = parseUtil.Int32ToPointer(app.InstanceCount)
		app.UpdateStatus = dao.ScaleFailed
	}

	//scale application
	if vebrType == "updateConfig" {
		app.UpdateStatus = dao.UpdateConfigFailed
		log.Errorf("input cpu %v", app.Cpu)
		log.Errorf("input memory %v", app.Memory)

		log.Errorf("rc before cpu %v", rc.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceCPU])
		log.Errorf("rc before memory %v", rc.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceMemory])
		rc.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceCPU] = resource.MustParse(app.Cpu)
		rc.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceMemory] = resource.MustParse(app.Memory)

		log.Errorf("rc late cpu %v", rc.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceCPU])
		log.Errorf("rc late memory %v", rc.Spec.Template.Spec.Containers[0].Resources.Requests[v1.ResourceMemory])
	}

	if vebrType == "redeployment" {
		//delete replicationController of before

		//delete conditions
		deleteOption := new(api.DeleteOptions)
		deleteOption.TypeMeta = unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}

		//if the value is 0 ,delete immediately. if not set,the default grace period for the specified type will be used ,
		//the default grace period is 30s
		deleteOption.GracePeriodSeconds = new(int64)

		//delete precondition(前提条件)
		deleteOption.Preconditions = &api.Preconditions{UID: &(rc.ObjectMeta.UID)}

		// If true/false，  added to/removed from the object's finalizers list
		deleteOption.OrphanDependents = parseUtil.BoolToPointer(false)
		err := dao.Clientset.Core().ReplicationControllers(NameSpace).Delete("name", deleteOption)
		if err != nil {
			log.Errorf("delete application failed ：%s", err.Error())
			code = StatusInternalServerError
			ret = map[string]interface{}{"success": false, "reason": err.Error()}
			return
		}

		//create a new replicationController
		rcTypeMeta := unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}

		rcObjectMeta := v1.ObjectMeta{
			Name: app.Name,
			Labels: map[string]string{
				"name": app.Name,
			},
		}

		rcSpec := v1.ReplicationControllerSpec{
			Replicas: parseUtil.Int32ToPointer(app.InstanceCount),
			Selector: map[string]string{
				"name": app.Name,
			},
			Template: &v1.PodTemplateSpec{
				v1.ObjectMeta{
					Name: app.Name,
					Labels: map[string]string{
						"name": app.Name,
					},
				},
				v1.PodSpec{
					RestartPolicy: v1.RestartPolicyAlways,
					NodeSelector: map[string]string{
						"name": app.Name,
					},
					Containers: []v1.Container{
						v1.Container{
							Name:  app.Name,
							Image: app.Image,
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: 6379,
									Protocol:      v1.ProtocolTCP,
								},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse(app.Cpu),
									v1.ResourceMemory: resource.MustParse(app.Memory),
								},
							},
						},
					},
				},
			},
		}

		rc := new(v1.ReplicationController)
		rc.TypeMeta = rcTypeMeta
		rc.ObjectMeta = rcObjectMeta
		rc.Spec = rcSpec
		app.UpdateStatus = dao.RedeploymentFailed
		result, err := dao.Clientset.Core().ReplicationControllers(NameSpace).Create(rc)
		if err != nil {
			log.Errorf("redeploy application failed ,the reason is %s", err.Error())
			if err = app.Update(); err != nil {
				log.Errorf("update application updateStatus failed,the reason is %s", err.Error())
			}
			code = StatusInternalServerError
			ret = JSON_EMPTY_OBJ
			return
		}

		app.UpdateStatus = dao.RedeploymentSuccessed
		app.Status = dao.AppRunning
		if err = app.Update(); err != nil {
			log.Errorf("update application updateStatus failed,the reason is %s", err.Error())
			code = StatusInternalServerError
			ret = JSON_EMPTY_OBJ
			return
		}

		code = StatusCreated
		ret = result
		return
	}

	return updateRc(NameSpace, app, rc)
}

func updateRc(nameSpace string, app *dao.App, rc *v1.ReplicationController) (code string, ret interface{}) {
	_, err := dao.Clientset.Core().ReplicationControllers(nameSpace).Update(rc)

	if err != nil {
		log.Errorf("update err :%v", err.Error())

		if err = app.Update(); err != nil {
			log.Errorf("update application failed,the reason is %s", err.Error())
			code = StatusInternalServerError
			ret = JSON_EMPTY_OBJ
			return
		}

		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	if err = app.Update(); err != nil {
		log.Errorf("update application failed,the reason is %s", err.Error())
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = OK
	return
}

func GetAppStatus(req *http.Request) (code string, ret interface{}) {
	NameSpace = req.FormValue("nameSpace")
	generateName := req.FormValue("appName") + "-"
	podList, err := dao.Clientset.Core().Pods(NameSpace).List(api.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	//get every container's status of pods
	for i := 0; i < len(podList.Items); i++ {
		if podList.Items[i].ObjectMeta.GenerateName == generateName && podList.Items[i].ObjectMeta.Namespace == NameSpace {
			app_status_slice = append(app_status_slice, podList.Items[i].Status.Phase)
		}
	}

	//determine the app status by all the container's status of pod
	for _, status := range app_status_slice {
		if status == v1.PodPending {
			app_status = status_Pending
			break
		}

		if status == v1.PodRunning {
			app_status = status_Running
		}

		if status == v1.PodSucceeded {
			app_status = status_Succeeded
		}

		if status == v1.PodFailed || status == v1.PodUnknown {
			app_status = status_Failed
			break
		}
	}

	code = StatusOK
	ret = app_status
	return
}

func GetAppContainers(req *http.Request) (code string, ret interface{}) {
	NameSpace = req.FormValue("nameSpace")
	generateName := req.FormValue("appName") + "-"

	podList, err := dao.Clientset.Core().Pods(NameSpace).List(api.ListOptions{})
	if err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_ARRAY
		log.Error(err.Error())
	}

	//get every container's status of pods
	var containers []dao.Container
	var pod v1.Pod
	for _, pod = range podList.Items {
		//Get the pod of app's
		if pod.ObjectMeta.GenerateName == generateName && pod.ObjectMeta.Namespace == NameSpace {
			container := dao.Container{}
			container.Name = pod.ObjectMeta.Name
			container.Image = pod.Spec.Containers[0].Image
			cpu := pod.Spec.Containers[0].Resources.Requests[v1.ResourceCPU]
			container.Cpu = cpu.String()
			memory := pod.Spec.Containers[0].Resources.Requests[v1.ResourceCPU]
			container.Memory = memory.String()
			container.Created = pod.Status.ContainerStatuses[0].State.Running.StartedAt.String()
			container.Ports = pod.Spec.Containers[0].Ports
			container.Envs = pod.Spec.Containers[0].Env
			container.IntranetIp = pod.Status.PodIP
			container.ExtranetIp = pod.Status.HostIP
			container.Mounts = pod.Spec.Containers[0].VolumeMounts

			containers = append(containers, container)
		}
	}

	code = StatusOK
	ret = containers
	return
}