class User {
    constructor(FullName, AvatarURL, Roles) {
        this.fullname = FullName
        this.avatarUrl = AvatarURL
        this.roles = Roles
    }

    hasRole(roleName) {
        // change later
        return this.roles.some(role => role.name === roleName);
    }
}

export {User}