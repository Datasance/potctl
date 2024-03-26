/*
 *  *******************************************************************************
 *  * Copyright (c) 2023 Datasance Teknoloji A.S.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package install

var pkg struct {
	scriptPrereq            string
	scriptInit              string
	scriptInstallDeps       string
	scriptInstallJava       string
	scriptInstallDocker     string
	scriptInstallIofog      string
	scriptUninstallIofog    string
	iofogDir                string
	agentDir                string
	scriptDatabaseMigration string
	scriptDatabaseSeeder    string
}

func init() {
	pkg.scriptPrereq = "check_prereqs.sh"
	pkg.scriptInit = "init.sh"
	pkg.scriptInstallDeps = "install_deps.sh"
	pkg.scriptInstallJava = "install_java.sh"
	pkg.scriptInstallDocker = "install_docker.sh"
	pkg.scriptInstallIofog = "install_iofog.sh"
	pkg.scriptUninstallIofog = "uninstall_iofog.sh"
	pkg.iofogDir = "/etc/iofog"
	pkg.agentDir = "/etc/iofog/agent"
	pkg.scriptDatabaseMigration = "db_migration_v1.0.0.sql"
	pkg.scriptDatabaseSeeder = "db_seeder_v1.0.0.sql"
}
