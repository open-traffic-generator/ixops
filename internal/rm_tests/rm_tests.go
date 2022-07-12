package rm_tests

import "github.com/open-traffic-generator/ixops/pkg/configs"

func RmTests(c configs.AppConfig) error {
	// home, _ := utils.GetHomeDirectory()
	// filePath := utils.ReturnPath([]string{home, c.IxOpsHome, c.IxiaC.Home, c.IxiaC.TestClient.Yaml})
	// _, err := utils.ExecCmd("kubectl", "delete", "-f", filePath)
	// if err != nil {
	// 	return err
	// }

	// _, err = utils.ExecCmd("kubectl", "delete", "namespace", "ixia-c-tests")
	// if err != nil {
	// 	return err
	// }

	// //Add wait for no namespace ixia-c-tests

	return nil
}
