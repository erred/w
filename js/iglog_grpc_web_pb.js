/**
 * @fileoverview gRPC-Web generated client stub for iglog
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.iglog = require('./iglog_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.iglog.FollowatchClient =
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
proto.iglog.FollowatchPromiseClient =
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
 *   !proto.iglog.Request,
 *   !proto.iglog.Events>}
 */
const methodInfo_Followatch_EventLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.iglog.Events,
  /** @param {!proto.iglog.Request} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.iglog.Events.deserializeBinary
);


/**
 * @param {!proto.iglog.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.iglog.Events)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.iglog.Events>|undefined}
 *     The XHR Node Readable Stream
 */
proto.iglog.FollowatchClient.prototype.eventLog =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/iglog.Followatch/EventLog',
      request,
      metadata || {},
      methodInfo_Followatch_EventLog,
      callback);
};


/**
 * @param {!proto.iglog.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.iglog.Events>}
 *     A native promise that resolves to the response
 */
proto.iglog.FollowatchPromiseClient.prototype.eventLog =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/iglog.Followatch/EventLog',
      request,
      metadata || {},
      methodInfo_Followatch_EventLog);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.iglog.Request,
 *   !proto.iglog.Users>}
 */
const methodInfo_Followatch_Followers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.iglog.Users,
  /** @param {!proto.iglog.Request} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.iglog.Users.deserializeBinary
);


/**
 * @param {!proto.iglog.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.iglog.Users)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.iglog.Users>|undefined}
 *     The XHR Node Readable Stream
 */
proto.iglog.FollowatchClient.prototype.followers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/iglog.Followatch/Followers',
      request,
      metadata || {},
      methodInfo_Followatch_Followers,
      callback);
};


/**
 * @param {!proto.iglog.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.iglog.Users>}
 *     A native promise that resolves to the response
 */
proto.iglog.FollowatchPromiseClient.prototype.followers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/iglog.Followatch/Followers',
      request,
      metadata || {},
      methodInfo_Followatch_Followers);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.iglog.Request,
 *   !proto.iglog.Users>}
 */
const methodInfo_Followatch_Following = new grpc.web.AbstractClientBase.MethodInfo(
  proto.iglog.Users,
  /** @param {!proto.iglog.Request} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.iglog.Users.deserializeBinary
);


/**
 * @param {!proto.iglog.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.iglog.Users)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.iglog.Users>|undefined}
 *     The XHR Node Readable Stream
 */
proto.iglog.FollowatchClient.prototype.following =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/iglog.Followatch/Following',
      request,
      metadata || {},
      methodInfo_Followatch_Following,
      callback);
};


/**
 * @param {!proto.iglog.Request} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.iglog.Users>}
 *     A native promise that resolves to the response
 */
proto.iglog.FollowatchPromiseClient.prototype.following =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/iglog.Followatch/Following',
      request,
      metadata || {},
      methodInfo_Followatch_Following);
};


module.exports = proto.iglog;

