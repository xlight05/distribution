/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package resources

import (
	"github.com/google/go-cmp/cmp"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestCreateGateway(t *testing.T) {
	tests := []struct {
		name string
		cell *v1alpha1.Cell
		want *v1alpha1.Gateway
	}{
		{
			name: "foo cell with empty gateway spec",
			cell: &v1alpha1.Cell{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
			},
			want: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo--gateway",
					Labels: map[string]string{
						vick.CellLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Cell",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
			},
		},
		{
			name: "foo cell with gateway spec",
			cell: &v1alpha1.Cell{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.CellSpec{
					GatewayTemplate: v1alpha1.GatewayTemplateSpec{
						Spec: v1alpha1.GatewaySpec{
							APIRoutes: []v1alpha1.APIRoute{
								{
									Context: "/context-1",
									Backend: "my-service",
									Global:  true,
									Definitions: []v1alpha1.APIDefinition{
										{
											Path:   "path1",
											Method: "GET",
										},
										{
											Path:   "path2",
											Method: "POST",
										},
									},
								},
							},
						},
					},
				},
			},
			want: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo--gateway",
					Labels: map[string]string{
						vick.CellLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Cell",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: v1alpha1.GatewaySpec{
					APIRoutes: []v1alpha1.APIRoute{
						{
							Context: "/context-1",
							Backend: "http://foo--my-service-service",
							Global:  true,
							Definitions: []v1alpha1.APIDefinition{
								{
									Path:   "path1",
									Method: "GET",
								},
								{
									Path:   "path2",
									Method: "POST",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateGateway(test.cell)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateGateway (-want, +got)\n%v", diff)
			}
		})
	}
}