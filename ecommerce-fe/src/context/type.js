class User {
    constructor(FullName, AvatarURL, Roles) {
        this.fullname = FullName;
        this.avatarUrl = AvatarURL;
        this.roles = Roles || [];
    }

    // Kiểm tra xem user có role cụ thể hay không
    hasRole(roleName) {
        if (!this.roles || this.roles.length === 0) return false;
        return this.roles.some(role => role.name === roleName);
    }

    // Kiểm tra xem user có bất kỳ role nào trong danh sách không
    hasAnyRole(roleNames) {
        if (!this.roles || this.roles.length === 0) return false;
        return roleNames.some(roleName => this.hasRole(roleName));
    }

    // Kiểm tra xem user có tất cả các role trong danh sách không
    hasAllRoles(roleNames) {
        if (!this.roles || this.roles.length === 0) return false;
        return roleNames.every(roleName => this.hasRole(roleName));
    }

    // Lấy danh sách tên các role của user
    getRoleNames() {
        if (!this.roles || this.roles.length === 0) return [];
        return this.roles.map(role => role.name);
    }
}

export { User };