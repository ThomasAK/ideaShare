# IdeaShare
A tool for sharing and voting on ideas

### Entities 
- Idea
  - ID
  - Title
  - Description
  - Status
  - CreateBy
  - CreatedTS
  - UpdatedTS
- IdeaComment
  - ID
  - IdeaID
  - Comment
  - CreatedBy
  - CreatedTS
  - UpdatedTS
- IdeaLike
  - ID
  - IdeaID
  - CreatedBy
  - CreatedTS
  - UpdatedTS
    - Unique: IdeaID,CreatedBy
- User
  - ID
  - ExternalID
  - FirstName
  - LastName
- UserRole
  - ID
  - UserID
  - Role
    - Roles: Approver, Admin