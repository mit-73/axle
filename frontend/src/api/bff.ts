import { createClient } from '@connectrpc/connect'
import { ProjectService } from '@axle/contracts/bff/v1/projects_pb'
import { UserService } from '@axle/contracts/bff/v1/users_pb'
import { transport } from './transport'

export const projectsClient = createClient(ProjectService, transport)
export const usersClient = createClient(UserService, transport)
