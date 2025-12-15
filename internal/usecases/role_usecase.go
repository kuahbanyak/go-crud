package usecases
import (
"context"
"errors"
"fmt"
"github.com/kuahbanyak/go-crud/internal/domain/entities"
"github.com/kuahbanyak/go-crud/internal/domain/repositories"
"github.com/kuahbanyak/go-crud/internal/shared/types"
"github.com/kuahbanyak/go-crud/pkg/pagination"
)
type RoleUsecase struct {
roleRepo repositories.RoleRepository
userRepo repositories.UserRepository
}
func NewRoleUsecase(roleRepo repositories.RoleRepository, userRepo repositories.UserRepository) *RoleUsecase {
return &RoleUsecase{
roleRepo: roleRepo,
userRepo: userRepo,
}
}
func (u *RoleUsecase) CreateRole(ctx context.Context, role *entities.Role) error {
existingRole, err := u.roleRepo.GetByName(ctx, role.Name)
if err != nil {
return fmt.Errorf("failed to check existing role: %w", err)
}
if existingRole != nil {
return errors.New("role with this name already exists")
}
return u.roleRepo.Create(ctx, role)
}
func (u *RoleUsecase) GetRoleByID(ctx context.Context, id types.MSSQLUUID) (*entities.Role, error) {
role, err := u.roleRepo.GetByID(ctx, id)
if err != nil {
return nil, err
}
if role == nil {
return nil, errors.New("role not found")
}
return role, nil
}
func (u *RoleUsecase) GetAllRoles(ctx context.Context) ([]*entities.Role, error) {
return u.roleRepo.GetAll(ctx)
}
func (u *RoleUsecase) GetAllRolesPaginated(ctx context.Context, pagParams pagination.Params, filterParams pagination.FilterParams) ([]*entities.Role, int64, error) {
return u.roleRepo.GetAllPaginated(ctx, pagParams, filterParams)
}
func (u *RoleUsecase) GetActiveRoles(ctx context.Context) ([]*entities.Role, error) {
return u.roleRepo.GetActive(ctx)
}
func (u *RoleUsecase) UpdateRole(ctx context.Context, id types.MSSQLUUID, updateData *entities.Role) (*entities.Role, error) {
existingRole, err := u.roleRepo.GetByID(ctx, id)
if err != nil {
return nil, err
}
if existingRole == nil {
return nil, errors.New("role not found")
}
if updateData.DisplayName != "" {
existingRole.DisplayName = updateData.DisplayName
}
if updateData.Description != "" {
existingRole.Description = updateData.Description
}
existingRole.IsActive = updateData.IsActive
err = u.roleRepo.Update(ctx, existingRole)
if err != nil {
return nil, err
}
return existingRole, nil
}
func (u *RoleUsecase) DeleteRole(ctx context.Context, id types.MSSQLUUID) error {
existingRole, err := u.roleRepo.GetByID(ctx, id)
if err != nil {
return err
}
if existingRole == nil {
return errors.New("role not found")
}
return u.roleRepo.Delete(ctx, id)
}
func (u *RoleUsecase) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy types.MSSQLUUID) error {
user, err := u.userRepo.GetByID(ctx, userID)
if err != nil {
return err
}
if user == nil {
return errors.New("user not found")
}
role, err := u.roleRepo.GetByID(ctx, roleID)
if err != nil {
return err
}
if role == nil {
return errors.New("role not found")
}
if !role.IsActive {
return errors.New("cannot assign inactive role")
}
return u.roleRepo.AssignRoleToUser(ctx, userID, roleID, assignedBy)
}
func (u *RoleUsecase) RemoveRoleFromUser(ctx context.Context, userID, roleID types.MSSQLUUID) error {
user, err := u.userRepo.GetByID(ctx, userID)
if err != nil {
return err
}
if user == nil {
return errors.New("user not found")
}
role, err := u.roleRepo.GetByID(ctx, roleID)
if err != nil {
return err
}
if role == nil {
return errors.New("role not found")
}
return u.roleRepo.RemoveRoleFromUser(ctx, userID, roleID)
}
func (u *RoleUsecase) GetUserRoles(ctx context.Context, userID types.MSSQLUUID) ([]*entities.Role, error) {
user, err := u.userRepo.GetByID(ctx, userID)
if err != nil {
return nil, err
}
if user == nil {
return nil, errors.New("user not found")
}
return u.roleRepo.GetUserRoles(ctx, userID)
}
func (u *RoleUsecase) HasRole(ctx context.Context, userID types.MSSQLUUID, roleName string) (bool, error) {
return u.roleRepo.HasRole(ctx, userID, roleName)
}
func (u *RoleUsecase) GetUsersByRole(ctx context.Context, roleID types.MSSQLUUID) ([]*entities.User, error) {
role, err := u.roleRepo.GetByID(ctx, roleID)
if err != nil {
return nil, err
}
if role == nil {
return nil, errors.New("role not found")
}
return u.roleRepo.GetUsersByRole(ctx, roleID)
}
