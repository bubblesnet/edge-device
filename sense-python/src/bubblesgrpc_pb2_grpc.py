# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import bubblesgrpc_pb2 as bubblesgrpc__pb2


class SensorStoreAndForwardStub(object):
  """The greeting service definition.
  """

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.StoreAndForward = channel.unary_unary(
        '/bubblesgrpc.SensorStoreAndForward/StoreAndForward',
        request_serializer=bubblesgrpc__pb2.SensorRequest.SerializeToString,
        response_deserializer=bubblesgrpc__pb2.SensorReply.FromString,
        )


class SensorStoreAndForwardServicer(object):
  """The greeting service definition.
  """

  def StoreAndForward(self, request, context):
    """Sends a greeting
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_SensorStoreAndForwardServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'StoreAndForward': grpc.unary_unary_rpc_method_handler(
          servicer.StoreAndForward,
          request_deserializer=bubblesgrpc__pb2.SensorRequest.FromString,
          response_serializer=bubblesgrpc__pb2.SensorReply.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'bubblesgrpc.SensorStoreAndForward', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
