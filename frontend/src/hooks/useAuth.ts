import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { authService, type LoginRequest, type RegisterRequest, type UpdateProfileRequest, type ChangePasswordRequest } from '@/services/authService';
import { useAuthStore } from '@/stores/authStore';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';

export function useLogin() {
  const { login } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (data: LoginRequest) => authService.login(data),
    onSuccess: (response) => {
      login(response.user, response.token);
      toast.success('登录成功');
      navigate('/');
    },
    onError: (error: Error) => {
      toast.error(error.message || '登录失败，请检查用户名和密码');
    },
  });
}

export function useRegister() {
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (data: RegisterRequest) => authService.register(data),
    onSuccess: () => {
      toast.success('注册成功');
      navigate('/login');
    },
    onError: (error: Error) => {
      toast.error(error.message || '注册失败');
    },
  });
}

export function useProfile() {
  const { isAuthenticated } = useAuthStore();

  return useQuery({
    queryKey: ['profile'],
    queryFn: () => authService.getProfile(),
    enabled: isAuthenticated,
  });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  const { setUser } = useAuthStore();

  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => authService.updateProfile(data),
    onSuccess: (user) => {
      setUser(user);
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      toast.success('个人信息更新成功');
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败');
    },
  });
}

export function useChangePassword() {
  return useMutation({
    mutationFn: (data: ChangePasswordRequest) => authService.changePassword(data),
    onSuccess: () => {
      toast.success('密码修改成功');
    },
    onError: (error: Error) => {
      toast.error(error.message || '密码修改失败');
    },
  });
}
