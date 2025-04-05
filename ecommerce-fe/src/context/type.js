class User {
    constructor(FullName, AvatarURL, Role) {
        this.fullname = FullName;
        this.avatarUrl = AvatarURL;
        this.role = Role;
    }

    // Kiểm tra xem user có role cụ thể hay không
    hasRole(roleName) {
        console.log(roleName)
        if (!this.role) return false;

        return this.role.name === roleName;
    }
}

export { User };