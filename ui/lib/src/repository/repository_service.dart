import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/api.dart';
import 'package:gitpods/src/api/api.dart' as api;
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/repository/repository_component.dart';
import 'package:gitpods/src/repository/repository_create_component.dart';

@Injectable()
class RepositoryService {
  RepositoryService(this._api);

  final API _api;

  Future<RepositoryPage> get(String owner, String name) async {
    api.Repository r = await _api.repositories.getRepository(owner, name);
    return RepositoryPage(Repository.fromAPI(r));
  }

  Future<Repository> create(
    String name,
    String description,
    String website,
  ) async {
    api.NewRepository newRepository = api.NewRepository();
    newRepository.name = name;
    newRepository.description = description;
    newRepository.website = website;

    try {
      api.Repository apiRepository =
          await _api.repositories.createRepository(newRepository);
      return Repository.fromAPI(apiRepository);
    } on api.ApiException catch (e) {
      if (e.code == 422) {
        throw new ValidationException.fromJSON(json.decode(e.message));
      } else {
        throw e;
      }
    }
  }
}
