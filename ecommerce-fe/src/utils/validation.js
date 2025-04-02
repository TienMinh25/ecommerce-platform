import * as yup from 'yup';

export const loginSchema = yup.object().shape({
  email: yup.string().email('Email không hợp lệ').required('Email là bắt buộc'),
  password: yup
    .string()
    .min(6, 'Mật khẩu phải có ít nhất 6 ký tự')
    .required('Mật khẩu là bắt buộc'),
});

export const registerSchema = yup.object().shape({
  full_name: yup
    .string()
    .required('Họ tên là bắt buộc')
    .min(2, 'Họ tên phải có ít nhất 2 ký tự'),
  email: yup.string().email('Email không hợp lệ').required('Email là bắt buộc'),
  password: yup
    .string()
    .min(6, 'Mật khẩu phải có ít nhất 6 ký tự')
    .required('Mật khẩu là bắt buộc'),
  agreeTerms: yup
    .boolean()
    .oneOf([true], 'Bạn phải đồng ý với điều khoản sử dụng')
    .required('Bạn phải đồng ý với điều khoản sử dụng'),
});
