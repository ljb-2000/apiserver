// Copyright © 2016 huang jia <huangjia@yfcloud.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"os"

	"github.com/op/go-logging"
)

func New(module string) *logging.Logger {
	log := logging.MustGetLogger(module)
	format := logging.MustStringFormatter(
		`%{color}%{time:2006/01/02 15:04:05} %{level:.5s} > %{message} [%{shortfile}] %{color:reset}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(-1, "")
	logging.SetBackend(backendLeveled, backendFormatter)
	return log
}
