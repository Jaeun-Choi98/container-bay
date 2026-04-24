class ApiClient {
  private baseURL: string | undefined;
  private defaultHeaders: HeadersInit

  constructor(baseURL: string | undefined = process.env.REACT_APP_API_BASE_URL) {
    this.baseURL = baseURL;
    this.defaultHeaders = {
      'Content-Type': 'application/json',
    }
  }

  private getHeaders(customHeaders?: HeadersInit): HeadersInit {
    const token = localStorage.getItem('authToken');
    const authHeader = token ? { Authorization: `Bearer ${token}` } : null;
    return {
      ...this.defaultHeaders,
      ...authHeader,
      ...customHeaders,
    };
  }

  async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;

    let response: Response;

    // 네트워크 에러만 catch
    try {
      response = await fetch(url, {
        ...options,
        headers: this.getHeaders(options?.headers),
      });
    } catch (e) {
      throw new Error(`network error: ${e}`);
    }

    // HTTP 에러는 별도 처리
    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(
        errorData?.message || `http status: ${response.status}`
      );
    }

    if (response.status === 204) {
      return {} as T;
    }

    return await response.json();
  }

  async get<T>(endpoint: string, hearders?: HeadersInit): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'GET',
      headers: hearders
    });
  }

  async post<T>(endpoint: string, data?: any, hearders?: HeadersInit): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
      headers: hearders
    });
  }

  async put<T>(endpoint: string, data?: any, hearders?: HeadersInit): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
      headers: hearders
    });
  }

  async patch<T>(endpoint: string, data?: any, hearders?: HeadersInit): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: JSON.stringify(data),
      headers: hearders
    });
  }

  async delete<T>(endpoint: string, hearders?: HeadersInit): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'DELETE',
      headers: hearders
    });
  }

  // 파일 업로드용
  async uploadFile<T>(endpoint: string, file: File, additionalData?: Record<string, string>): Promise<T> {
    const formData = new FormData();
    formData.append('file', file);
    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value);
      });
    }

    const token = localStorage.getItem('authToken');
    const authHeader = token ? { Authorization: `Bearer ${token}` } : undefined;

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'POST',
      headers: authHeader,
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload failed! status: ${response.status}`);
    }

    return await response.json();
  }

}

export const apiClient = new ApiClient();