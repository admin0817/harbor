# coding: utf-8

"""
    Harbor API

    These APIs provide services for manipulating Harbor project. 

    OpenAPI spec version: 1.4.0
    
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""


from __future__ import absolute_import

import os
import sys
sys.path.append(os.environ["SWAGGER_CLIENT_PATH"])

import unittest
import testutils
import swagger_client
from swagger_client.models.project_req import ProjectReq
from swagger_client.models.configurations import Configurations
from swagger_client.rest import ApiException
from swagger_client.models.configurations import Configurations 
from pprint import pprint


#Testcase
# Define a LDAP group with harbor admin
class TestLdapAdminRole(unittest.TestCase):
    """AccessLog unit test stubs"""
    product_api = testutils.GetProductApi("admin", "Harbor12345")
    mike_product_api = testutils.GetProductApi("mike", "zhu88jie")
    project_id = 0
    def setUp(self):
        pass

    def tearDown(self):
        if self.project_id > 0 :
            self.mike_product_api.projects_project_id_delete(project_id=self.project_id)
        pass

    def testLdapAdminRole(self):
        """Test LdapAdminRole"""
        result = self.product_api.configurations_put(configurations=Configurations(ldap_group_admin_dn="cn=harbor_users,ou=groups,dc=example,dc=com"))
        pprint(result)

        # Create a private project
        result = self.product_api.projects_post(project=ProjectReq(project_name="test_private"))
        pprint(result)

        # query project with ldap user mike
        projects = self.mike_product_api.projects_get(name="test_private")
        self.assertTrue(projects.count>1)
        self.project_id = projects[0].project_id

        # check the mike is not admin in Database
        user_list = self.product_api.users_get(username="mike")
        pprint(user_list[0])
        self.assertFalse(user_list[0].sysadmin_flag)

        pass


if __name__ == '__main__':
    unittest.main()
