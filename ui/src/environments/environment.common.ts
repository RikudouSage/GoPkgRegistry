function getWindowEnvOrDefault(name: string, defaultValue: string): string {
  const anyWindow = window as any;

  if (typeof anyWindow.runtimeVariables === 'undefined') {
    return defaultValue;
  }

  return anyWindow.runtimeVariables[name] ?? defaultValue;
}

function getNodeEnvOrDefault(name: string, defaultValue: string): string {
  const val = (globalThis as any)?.process?.env?.[name] as string | undefined;
  return val ?? defaultValue;
}

export function getEnvOrDefault(name: string, defaultValue: string): string {
  if (typeof window === 'undefined') {
    return getNodeEnvOrDefault(name, defaultValue);
  }

  return getWindowEnvOrDefault(name, defaultValue);
}

export interface Environment {
  apiUrl: string;
  sameOrigin: boolean;
}

export const commonEnvironment = {
  apiUrl: getEnvOrDefault('API_URL', '/api'),
  sameOrigin: getEnvOrDefault('SAME_ORIGIN', 'true') === 'true',
};
