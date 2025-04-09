// User.js
class User {
    constructor(fullName, avatarUrl, roles) {
        this.fullname = fullName;
        this.avatarUrl = avatarUrl || '';
        this.roles = roles; // Thay đổi từ role sang roles (array)
    }

    // Getter cho role chính (đầu tiên trong danh sách)
    get primaryRole() {
        return this.roles && this.roles.length > 0 ? this.roles[0] : null;
    }

    // Kiểm tra user có role nào đó không
    hasRole(roleName) {
        if (!this.roles || this.roles.length === 0) return false;
        return this.roles.some(role => role.name === roleName);
    }
}

export default User;

export { User };