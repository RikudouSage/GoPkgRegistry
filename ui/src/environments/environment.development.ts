import { commonEnvironment, Environment } from './environment.common';

export const environment: Environment = {
  ...commonEnvironment,
  sameOrigin: false,
  apiUrl: 'http://localhost:8080',
};
