/**
 * @fileoverview gRPC-Web generated client stub for authed
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.authed = require('./authed_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.authed.authedClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.authed.authedPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.authed.Msg,
 *   !proto.authed.Msg>}
 */
const methodInfo_authed_Echo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.authed.Msg,
  /** @param {!proto.authed.Msg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.authed.Msg.deserializeBinary
);


/**
 * @param {!proto.authed.Msg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.authed.Msg)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.authed.Msg>|undefined}
 *     The XHR Node Readable Stream
 */
proto.authed.authedClient.prototype.echo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/authed.authed/Echo',
      request,
      metadata || {},
      methodInfo_authed_Echo,
      callback);
};


/**
 * @param {!proto.authed.Msg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.authed.Msg>}
 *     A native promise that resolves to the response
 */
proto.authed.authedPromiseClient.prototype.echo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/authed.authed/Echo',
      request,
      metadata || {},
      methodInfo_authed_Echo);
};


module.exports = proto.authed;

