ۺ
��
8
Const
output"dtype"
valuetensor"
dtypetype

NoOp
�
PartitionedCall
args2Tin
output2Tout"
Tin
list(type)("
Tout
list(type)("	
ffunc"
configstring "
config_protostring "
executor_typestring 
C
Placeholder
output"dtype"
dtypetype"
shapeshape:
�
StatefulPartitionedCall
args2Tin
output2Tout"
Tin
list(type)("
Tout
list(type)("	
ffunc"
configstring "
config_protostring "
executor_typestring �"serve*1.15.32v1.15.2-30-g4386a668�g

NoOpNoOp
�	
ConstConst"/device:CPU:0*
_output_shapes
: *
dtype0*�
value�B� B�
�
layer-0
layer-1
layer-2
	variables
trainable_variables
regularization_losses
	keras_api

signatures
R
		variables

trainable_variables
regularization_losses
	keras_api
R
	variables
trainable_variables
regularization_losses
	keras_api
h
_callable_losses
	variables
trainable_variables
regularization_losses
	keras_api
 
 
 
�
layer_regularization_losses

layers
	variables
metrics
trainable_variables
regularization_losses
non_trainable_variables
 
 
 
 
�
layer_regularization_losses

layers
		variables
metrics

trainable_variables
regularization_losses
non_trainable_variables
 
 
 
�
layer_regularization_losses

layers
	variables
 metrics
trainable_variables
regularization_losses
!non_trainable_variables
 
 
 
 
�
"layer_regularization_losses

#layers
	variables
$metrics
trainable_variables
regularization_losses
%non_trainable_variables
 

0
1
2
 
 
 
 
 
 
 
 
 
 
 
 
 
 
l
serving_default_xPlaceholder*#
_output_shapes
:���������*
dtype0*
shape:���������
l
serving_default_yPlaceholder*#
_output_shapes
:���������*
shape:���������*
dtype0
�
PartitionedCallPartitionedCallserving_default_xserving_default_y*
Tin
2**
config_proto

GPU 

CPU2J 8**
_gradient_op_typePartitionedCall-136*)
f$R"
 __inference_signature_wrapper_92*#
_output_shapes
:���������*
Tout
2
O
saver_filenamePlaceholder*
dtype0*
_output_shapes
: *
shape: 
�
StatefulPartitionedCallStatefulPartitionedCallsaver_filenameConst**
config_proto

GPU 

CPU2J 8*%
f R
__inference__traced_save_153*
Tin
2*
_output_shapes
: **
_gradient_op_typePartitionedCall-154*
Tout
2
�
StatefulPartitionedCall_1StatefulPartitionedCallsaver_filename*
_output_shapes
: *
Tin
2*(
f#R!
__inference__traced_restore_164*
Tout
2**
_gradient_op_typePartitionedCall-165**
config_proto

GPU 

CPU2J 8�U
�
;
__inference__wrapped_model_14
x
y
identityJ
model/out/addAddV2xy*#
_output_shapes
:���������*
T0U
IdentityIdentitymodel/out/add:z:0*
T0*#
_output_shapes
:���������"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:! 

_user_specified_namex:!

_user_specified_namey
�
h
<__inference_out_layer_call_and_return_conditional_losses_125
inputs_0
inputs_1
identityN
addAddV2inputs_0inputs_1*#
_output_shapes
:���������*
T0K
IdentityIdentityadd:z:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:( $
"
_user_specified_name
inputs/0:($
"
_user_specified_name
inputs/1
�
r
__inference__traced_save_153
file_prefix
savev2_const

identity_1��MergeV2Checkpoints�SaveV2�
StringJoin/inputs_1Const"/device:CPU:0*
dtype0*
_output_shapes
: *<
value3B1 B+_temp_11653721dfca4e5ab66e9510d72ef27e/parts

StringJoin
StringJoinfile_prefixStringJoin/inputs_1:output:0"/device:CPU:0*
N*
_output_shapes
: L

num_shardsConst*
dtype0*
value	B :*
_output_shapes
: f
ShardedFilename/shardConst"/device:CPU:0*
value	B : *
_output_shapes
: *
dtype0�
ShardedFilenameShardedFilenameStringJoin:output:0ShardedFilename/shard:output:0num_shards:output:0"/device:CPU:0*
_output_shapes
: �
SaveV2/tensor_namesConst"/device:CPU:0*1
value(B&B_CHECKPOINTABLE_OBJECT_GRAPH*
_output_shapes
:*
dtype0o
SaveV2/shape_and_slicesConst"/device:CPU:0*
valueB
B *
dtype0*
_output_shapes
:�
SaveV2SaveV2ShardedFilename:filename:0SaveV2/tensor_names:output:0 SaveV2/shape_and_slices:output:0savev2_const"/device:CPU:0*
_output_shapes
 *
dtypes
2�
&MergeV2Checkpoints/checkpoint_prefixesPackShardedFilename:filename:0^SaveV2"/device:CPU:0*
_output_shapes
:*
N*
T0�
MergeV2CheckpointsMergeV2Checkpoints/MergeV2Checkpoints/checkpoint_prefixes:output:0file_prefix^SaveV2"/device:CPU:0*
_output_shapes
 f
IdentityIdentityfile_prefix^MergeV2Checkpoints"/device:CPU:0*
_output_shapes
: *
T0h

Identity_1IdentityIdentity:output:0^MergeV2Checkpoints^SaveV2*
T0*
_output_shapes
: "!

identity_1Identity_1:output:0*
_input_shapes
: : 2
SaveV2SaveV22(
MergeV2CheckpointsMergeV2Checkpoints:+ '
%
_user_specified_namefile_prefix: 
�
e
;__inference_out_layer_call_and_return_conditional_losses_28

inputs
inputs_1
identityL
addAddV2inputsinputs_1*
T0*#
_output_shapes
:���������K
IdentityIdentityadd:z:0*
T0*#
_output_shapes
:���������"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:& "
 
_user_specified_nameinputs:&"
 
_user_specified_nameinputs
�
@
"__inference_model_layer_call_fn_68
x
y
identity�
PartitionedCallPartitionedCallxy*)
_gradient_op_typePartitionedCall-65*
Tin
2*F
fAR?
=__inference_model_layer_call_and_return_conditional_losses_64*
Tout
2*#
_output_shapes
:���������**
config_proto

GPU 

CPU2J 8\
IdentityIdentityPartitionedCall:output:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:! 

_user_specified_namex:!

_user_specified_namey
�
[
=__inference_model_layer_call_and_return_conditional_losses_45
x
y
identity�
out/PartitionedCallPartitionedCallxy*#
_output_shapes
:���������*
Tout
2**
config_proto

GPU 

CPU2J 8*D
f?R=
;__inference_out_layer_call_and_return_conditional_losses_28*)
_gradient_op_typePartitionedCall-36*
Tin
2`
IdentityIdentityout/PartitionedCall:output:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:! 

_user_specified_namex:!

_user_specified_namey
�
M
!__inference_out_layer_call_fn_131
inputs_0
inputs_1
identity�
PartitionedCallPartitionedCallinputs_0inputs_1*#
_output_shapes
:���������*
Tout
2*D
f?R=
;__inference_out_layer_call_and_return_conditional_losses_28*
Tin
2*)
_gradient_op_typePartitionedCall-36**
config_proto

GPU 

CPU2J 8\
IdentityIdentityPartitionedCall:output:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:( $
"
_user_specified_name
inputs/0:($
"
_user_specified_name
inputs/1
�
g
=__inference_model_layer_call_and_return_conditional_losses_64

inputs
inputs_1
identity�
out/PartitionedCallPartitionedCallinputsinputs_1*
Tout
2**
config_proto

GPU 

CPU2J 8*D
f?R=
;__inference_out_layer_call_and_return_conditional_losses_28*)
_gradient_op_typePartitionedCall-36*
Tin
2*#
_output_shapes
:���������`
IdentityIdentityout/PartitionedCall:output:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:& "
 
_user_specified_nameinputs:&"
 
_user_specified_nameinputs
�
@
"__inference_model_layer_call_fn_84
x
y
identity�
PartitionedCallPartitionedCallxy**
config_proto

GPU 

CPU2J 8*
Tout
2*)
_gradient_op_typePartitionedCall-81*#
_output_shapes
:���������*
Tin
2*F
fAR?
=__inference_model_layer_call_and_return_conditional_losses_80\
IdentityIdentityPartitionedCall:output:0*
T0*#
_output_shapes
:���������"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:! 

_user_specified_namex:!

_user_specified_namey
�
j
>__inference_model_layer_call_and_return_conditional_losses_101
inputs_0
inputs_1
identityR
out/addAddV2inputs_0inputs_1*#
_output_shapes
:���������*
T0O
IdentityIdentityout/add:z:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:( $
"
_user_specified_name
inputs/0:($
"
_user_specified_name
inputs/1
�
O
#__inference_model_layer_call_fn_113
inputs_0
inputs_1
identity�
PartitionedCallPartitionedCallinputs_0inputs_1**
config_proto

GPU 

CPU2J 8*)
_gradient_op_typePartitionedCall-65*
Tin
2*F
fAR?
=__inference_model_layer_call_and_return_conditional_losses_64*
Tout
2*#
_output_shapes
:���������\
IdentityIdentityPartitionedCall:output:0*
T0*#
_output_shapes
:���������"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:( $
"
_user_specified_name
inputs/0:($
"
_user_specified_name
inputs/1
�
[
=__inference_model_layer_call_and_return_conditional_losses_54
x
y
identity�
out/PartitionedCallPartitionedCallxy*
Tin
2**
config_proto

GPU 

CPU2J 8*D
f?R=
;__inference_out_layer_call_and_return_conditional_losses_28*
Tout
2*)
_gradient_op_typePartitionedCall-36*#
_output_shapes
:���������`
IdentityIdentityout/PartitionedCall:output:0*
T0*#
_output_shapes
:���������"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:! 

_user_specified_namex:!

_user_specified_namey
�
>
 __inference_signature_wrapper_92
x
y
identity�
PartitionedCallPartitionedCallxy*
Tout
2*
Tin
2*)
_gradient_op_typePartitionedCall-89*#
_output_shapes
:���������*&
f!R
__inference__wrapped_model_14**
config_proto

GPU 

CPU2J 8\
IdentityIdentityPartitionedCall:output:0*
T0*#
_output_shapes
:���������"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:! 

_user_specified_namex:!

_user_specified_namey
�
Q
__inference__traced_restore_164
file_prefix

identity_1��	RestoreV2�
RestoreV2/tensor_namesConst"/device:CPU:0*
_output_shapes
:*
dtype0*1
value(B&B_CHECKPOINTABLE_OBJECT_GRAPHr
RestoreV2/shape_and_slicesConst"/device:CPU:0*
dtype0*
valueB
B *
_output_shapes
:�
	RestoreV2	RestoreV2file_prefixRestoreV2/tensor_names:output:0#RestoreV2/shape_and_slices:output:0"/device:CPU:0*
_output_shapes
:*
dtypes
21
NoOpNoOp"/device:CPU:0*
_output_shapes
 X
IdentityIdentityfile_prefix^NoOp"/device:CPU:0*
T0*
_output_shapes
: V

Identity_1IdentityIdentity:output:0
^RestoreV2*
T0*
_output_shapes
: "!

identity_1Identity_1:output:0*
_input_shapes
: 2
	RestoreV2	RestoreV2:+ '
%
_user_specified_namefile_prefix
�
O
#__inference_model_layer_call_fn_119
inputs_0
inputs_1
identity�
PartitionedCallPartitionedCallinputs_0inputs_1*#
_output_shapes
:���������*)
_gradient_op_typePartitionedCall-81*
Tin
2*F
fAR?
=__inference_model_layer_call_and_return_conditional_losses_80*
Tout
2**
config_proto

GPU 

CPU2J 8\
IdentityIdentityPartitionedCall:output:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:( $
"
_user_specified_name
inputs/0:($
"
_user_specified_name
inputs/1
�
j
>__inference_model_layer_call_and_return_conditional_losses_107
inputs_0
inputs_1
identityR
out/addAddV2inputs_0inputs_1*#
_output_shapes
:���������*
T0O
IdentityIdentityout/add:z:0*#
_output_shapes
:���������*
T0"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:( $
"
_user_specified_name
inputs/0:($
"
_user_specified_name
inputs/1
�
g
=__inference_model_layer_call_and_return_conditional_losses_80

inputs
inputs_1
identity�
out/PartitionedCallPartitionedCallinputsinputs_1*#
_output_shapes
:���������*
Tout
2*D
f?R=
;__inference_out_layer_call_and_return_conditional_losses_28*
Tin
2*)
_gradient_op_typePartitionedCall-36**
config_proto

GPU 

CPU2J 8`
IdentityIdentityout/PartitionedCall:output:0*
T0*#
_output_shapes
:���������"
identityIdentity:output:0*1
_input_shapes 
:���������:���������:& "
 
_user_specified_nameinputs:&"
 
_user_specified_nameinputs"�J
saver_filename:0StatefulPartitionedCall:0StatefulPartitionedCall_18"
saved_model_main_op

NoOp*�
serving_default�
+
x&
serving_default_x:0���������
+
y&
serving_default_y:0���������+
out$
PartitionedCall:0���������tensorflow/serving/predict*>
__saved_model_init_op%#
__saved_model_init_op

NoOp:�K
�
layer-0
layer-1
layer-2
	variables
trainable_variables
regularization_losses
	keras_api

signatures
&__call__
'_default_save_signature
*(&call_and_return_all_conditional_losses"�
_tf_keras_model�{"class_name": "Model", "name": "model", "trainable": true, "expects_training_arg": true, "dtype": null, "batch_input_shape": null, "config": {"name": "model", "layers": [{"name": "x", "class_name": "InputLayer", "config": {"batch_input_shape": [null], "dtype": "float32", "sparse": false, "ragged": false, "name": "x"}, "inbound_nodes": []}, {"name": "y", "class_name": "InputLayer", "config": {"batch_input_shape": [null], "dtype": "float32", "sparse": false, "ragged": false, "name": "y"}, "inbound_nodes": []}, {"name": "out", "class_name": "Add", "config": {"name": "out", "trainable": true, "dtype": "float32"}, "inbound_nodes": [[["x", 0, 0, {}], ["y", 0, 0, {}]]]}], "input_layers": [["x", 0, 0], ["y", 0, 0]], "output_layers": [["out", 0, 0]]}, "input_spec": [null, null], "activity_regularizer": null, "keras_version": "2.2.4-tf", "backend": "tensorflow", "model_config": {"class_name": "Model", "config": {"name": "model", "layers": [{"name": "x", "class_name": "InputLayer", "config": {"batch_input_shape": [null], "dtype": "float32", "sparse": false, "ragged": false, "name": "x"}, "inbound_nodes": []}, {"name": "y", "class_name": "InputLayer", "config": {"batch_input_shape": [null], "dtype": "float32", "sparse": false, "ragged": false, "name": "y"}, "inbound_nodes": []}, {"name": "out", "class_name": "Add", "config": {"name": "out", "trainable": true, "dtype": "float32"}, "inbound_nodes": [[["x", 0, 0, {}], ["y", 0, 0, {}]]]}], "input_layers": [["x", 0, 0], ["y", 0, 0]], "output_layers": [["out", 0, 0]]}}}
�
		variables

trainable_variables
regularization_losses
	keras_api
)__call__
**&call_and_return_all_conditional_losses"�
_tf_keras_layer�{"class_name": "InputLayer", "name": "x", "trainable": true, "expects_training_arg": true, "dtype": "float32", "batch_input_shape": [null], "config": {"batch_input_shape": [null], "dtype": "float32", "sparse": false, "ragged": false, "name": "x"}, "input_spec": null, "activity_regularizer": null}
�
	variables
trainable_variables
regularization_losses
	keras_api
+__call__
*,&call_and_return_all_conditional_losses"�
_tf_keras_layer�{"class_name": "InputLayer", "name": "y", "trainable": true, "expects_training_arg": true, "dtype": "float32", "batch_input_shape": [null], "config": {"batch_input_shape": [null], "dtype": "float32", "sparse": false, "ragged": false, "name": "y"}, "input_spec": null, "activity_regularizer": null}
�
_callable_losses
	variables
trainable_variables
regularization_losses
	keras_api
-__call__
*.&call_and_return_all_conditional_losses"�
_tf_keras_layer�{"class_name": "Add", "name": "out", "trainable": true, "expects_training_arg": false, "dtype": "float32", "batch_input_shape": null, "config": {"name": "out", "trainable": true, "dtype": "float32"}, "input_spec": null, "activity_regularizer": null}
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
�
layer_regularization_losses

layers
	variables
metrics
trainable_variables
regularization_losses
non_trainable_variables
&__call__
'_default_save_signature
*(&call_and_return_all_conditional_losses
&("call_and_return_conditional_losses"
_generic_user_object
,
/serving_default"
signature_map
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
�
layer_regularization_losses

layers
		variables
metrics

trainable_variables
regularization_losses
non_trainable_variables
)__call__
**&call_and_return_all_conditional_losses
&*"call_and_return_conditional_losses"
_generic_user_object
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
�
layer_regularization_losses

layers
	variables
 metrics
trainable_variables
regularization_losses
!non_trainable_variables
+__call__
*,&call_and_return_all_conditional_losses
&,"call_and_return_conditional_losses"
_generic_user_object
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
�
"layer_regularization_losses

#layers
	variables
$metrics
trainable_variables
regularization_losses
%non_trainable_variables
-__call__
*.&call_and_return_all_conditional_losses
&."call_and_return_conditional_losses"
_generic_user_object
 "
trackable_list_wrapper
5
0
1
2"
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
�2�
#__inference_model_layer_call_fn_119
#__inference_model_layer_call_fn_113
"__inference_model_layer_call_fn_68
"__inference_model_layer_call_fn_84�
���
FullArgSpec1
args)�&
jself
jinputs

jtraining
jmask
varargs
 
varkw
 
defaults�
p 

 

kwonlyargs� 
kwonlydefaults� 
annotations� *
 
�2�
__inference__wrapped_model_14�
���
FullArgSpec
args� 
varargsjargs
varkw
 
defaults
 

kwonlyargs� 
kwonlydefaults
 
annotations� *:�7
5�2
�
x���������
�
y���������
�2�
>__inference_model_layer_call_and_return_conditional_losses_101
=__inference_model_layer_call_and_return_conditional_losses_45
>__inference_model_layer_call_and_return_conditional_losses_107
=__inference_model_layer_call_and_return_conditional_losses_54�
���
FullArgSpec1
args)�&
jself
jinputs

jtraining
jmask
varargs
 
varkw
 
defaults�
p 

 

kwonlyargs� 
kwonlydefaults� 
annotations� *
 
�2��
���
FullArgSpec
args�
jself
jinputs
varargs
 
varkwjkwargs
defaults� 

kwonlyargs�

jtraining%
kwonlydefaults�

trainingp 
annotations� *
 
�2��
���
FullArgSpec
args�
jself
jinputs
varargs
 
varkwjkwargs
defaults� 

kwonlyargs�

jtraining%
kwonlydefaults�

trainingp 
annotations� *
 
�2��
���
FullArgSpec
args�
jself
jinputs
varargs
 
varkwjkwargs
defaults� 

kwonlyargs�

jtraining%
kwonlydefaults�

trainingp 
annotations� *
 
�2��
���
FullArgSpec
args�
jself
jinputs
varargs
 
varkwjkwargs
defaults� 

kwonlyargs�

jtraining%
kwonlydefaults�

trainingp 
annotations� *
 
�2�
!__inference_out_layer_call_fn_131�
���
FullArgSpec
args�
jself
jinputs
varargs
 
varkw
 
defaults
 

kwonlyargs� 
kwonlydefaults
 
annotations� *
 
�2�
<__inference_out_layer_call_and_return_conditional_losses_125�
���
FullArgSpec
args�
jself
jinputs
varargs
 
varkw
 
defaults
 

kwonlyargs� 
kwonlydefaults
 
annotations� *
 
*B(
 __inference_signature_wrapper_92xy�
=__inference_model_layer_call_and_return_conditional_losses_54qL�I
B�?
5�2
�
x���������
�
y���������
p 

 
� "!�
�
0���������
� �
!__inference_out_layer_call_fn_131jR�O
H�E
C�@
�
inputs/0���������
�
inputs/1���������
� "�����������
#__inference_model_layer_call_fn_119rZ�W
P�M
C�@
�
inputs/0���������
�
inputs/1���������
p 

 
� "�����������
>__inference_model_layer_call_and_return_conditional_losses_101Z�W
P�M
C�@
�
inputs/0���������
�
inputs/1���������
p

 
� "!�
�
0���������
� �
 __inference_signature_wrapper_92rI�F
� 
?�<

x�
x���������

y�
y���������"%�"
 
out�
out����������
"__inference_model_layer_call_fn_84dL�I
B�?
5�2
�
x���������
�
y���������
p 

 
� "�����������
__inference__wrapped_model_14mD�A
:�7
5�2
�
x���������
�
y���������
� "%�"
 
out�
out����������
=__inference_model_layer_call_and_return_conditional_losses_45qL�I
B�?
5�2
�
x���������
�
y���������
p

 
� "!�
�
0���������
� �
#__inference_model_layer_call_fn_113rZ�W
P�M
C�@
�
inputs/0���������
�
inputs/1���������
p

 
� "�����������
>__inference_model_layer_call_and_return_conditional_losses_107Z�W
P�M
C�@
�
inputs/0���������
�
inputs/1���������
p 

 
� "!�
�
0���������
� �
"__inference_model_layer_call_fn_68dL�I
B�?
5�2
�
x���������
�
y���������
p

 
� "�����������
<__inference_out_layer_call_and_return_conditional_losses_125wR�O
H�E
C�@
�
inputs/0���������
�
inputs/1���������
� "!�
�
0���������
� 