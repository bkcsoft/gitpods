import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/user.dart';
import 'package:http/browser_client.dart';
import 'package:http/http.dart';

const userMe = '''
query me {
  me {
    id
    email
    username
    name
    created_at
    updated_at
  }
}
''';

const usersQuery = '''
query UsersQuery {
  users {
    id
    email
    username
    name
    created_at
    updated_at
  }
}
''';

const userProfile = '''
query userProfile(\$username: String!) {
  user(username: \$username) {
    id
    username
    name
    email
    created_at
    updated_at
    repositories {
      id
      name
      description
      forks
      stars
    }
  }
}
''';

const userUpdate = '''
mutation updateUser(\$id: ID!, \$user: updatedUser!) {
  updateUser(id: \$id, user: \$user) {
    id
    email
    username
    name
    created_at
    updated_at
  }
}
''';

@Injectable()
class UserService {
  final BrowserClient _http;

  UserService(this._http);

  Future<User> me() async {
    var payload = JSON.encode({
      'query': userMe,
    });

    Response resp = await this._http.post('/api/query', body: payload);

    if (resp.statusCode != 200) {
      var body = JSON.decode(resp.body);
      throw new UnauthorizedError(body['errors'][0]['detail']);
    }

    var body = JSON.decode(resp.body);
    return new User.fromJSON(body['data']['me']);
  }

  Future<List<User>> list() async {
    var payload = JSON.encode({
      'query': usersQuery,
    });

    Response resp = await this._http.post('/api/query', body: payload);

    var body = JSON.decode(resp.body);
    return body['data']['users']
        .map((json) => new User.fromJSON(json))
        .toList();
  }

  Future<User> profile(String username) async {
    var payload = JSON.encode({
      'query': userProfile,
      'variables': {
        'username': username,
      }
    });

    Response resp = await this._http.post('/api/query', body: payload);

    var body = JSON.decode(resp.body);
    return new User.fromJSON(body['data']['user']);
  }

  Future<User> update(User user) async {
    var payload = JSON.encode({
      'query': userUpdate,
      'variables': {
        'id': user.id,
        'user': {
          'name': user.name,
        },
      },
    });

    Response resp = await this._http.post('/api/query', body: payload);

    var body = JSON.decode(resp.body);
    return new User.fromJSON(body['data']['updateUser']);
  }
}

class UnauthorizedError extends Error {
  final String message;

  UnauthorizedError(this.message);

  @override
  String toString() => "Unauthorizted state: $message";
}
