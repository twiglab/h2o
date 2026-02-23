<template>
  <div class="tags-view-container">
    <el-scrollbar ref="scrollbarRef" class="tags-view-wrapper">
      <router-link
        v-for="tag in visitedViews"
        :key="tag.path"
        :to="tag.fullPath"
        :class="['tags-view-item', { active: isActive(tag) }]"
        @click.middle="!isAffix(tag) && closeTag(tag)"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        {{ tag.title }}
        <el-icon
          v-if="!isAffix(tag)"
          class="close-icon"
          @click.prevent.stop="closeTag(tag)"
        >
          <Close />
        </el-icon>
      </router-link>
    </el-scrollbar>

    <!-- 右键菜单 -->
    <ul v-show="menuVisible" :style="{ left: menuLeft + 'px', top: menuTop + 'px' }" class="context-menu">
      <li @click="refreshSelectedTag(selectedTag)">
        <el-icon><Refresh /></el-icon>
        刷新
      </li>
      <li v-if="!isAffix(selectedTag)" @click="closeTag(selectedTag)">
        <el-icon><Close /></el-icon>
        关闭
      </li>
      <li @click="closeOthersTags">
        <el-icon><CircleClose /></el-icon>
        关闭其他
      </li>
      <li @click="closeAllTags">
        <el-icon><FolderDelete /></el-icon>
        关闭所有
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { RouteLocationNormalized } from 'vue-router'
import { useTagsViewStore, usePermissionStore } from '@/stores'
import type { TagView } from '@/stores/modules/tagsView'

const route = useRoute()
const router = useRouter()
const tagsViewStore = useTagsViewStore()
const permissionStore = usePermissionStore()

const scrollbarRef = ref()

// 已访问的标签
const visitedViews = computed(() => tagsViewStore.visitedViews)

// 判断是否为当前激活
function isActive(tag: TagView): boolean {
  return tag.path === route.path
}

// 判断是否为固定标签
function isAffix(tag: TagView | null): boolean {
  return tag?.affix ?? false
}

// 添加标签
function addTags() {
  if (route.name && route.meta?.title) {
    tagsViewStore.addView(route as RouteLocationNormalized)
  }
}

// 关闭标签
function closeTag(view: TagView | null) {
  if (!view) return
  tagsViewStore.removeView(view)
  if (isActive(view)) {
    toLastView()
  }
}

// 跳转到最后一个标签
function toLastView() {
  const latestView = visitedViews.value.slice(-1)[0]
  if (latestView) {
    router.push(latestView.path)
  } else {
    router.push('/')
  }
}

// 关闭其他标签
function closeOthersTags() {
  if (selectedTag.value) {
    tagsViewStore.removeOtherViews(selectedTag.value)
    if (!isActive(selectedTag.value)) {
      router.push(selectedTag.value.path)
    }
  }
  closeMenu()
}

// 关闭所有标签
function closeAllTags() {
  tagsViewStore.removeAllViews()
  toLastView()
  closeMenu()
}

// 刷新选中的标签
function refreshSelectedTag(view: TagView | null) {
  if (!view) return
  tagsViewStore.removeCachedView(view)
  nextTick(() => {
    router.replace('/redirect' + view.fullPath)
  })
  closeMenu()
}

// 初始化固定标签
function initAffixTags() {
  const routes = permissionStore.routes
  const filterAffixTags = (routes: any[], basePath = '/'): TagView[] => {
    let tags: TagView[] = []
    routes.forEach(route => {
      if (route.meta?.affix) {
        const tagPath = basePath + (route.path.startsWith('/') ? route.path : '/' + route.path)
        tags.push({
          name: route.name,
          title: route.meta.title,
          path: tagPath.replace('//', '/'),
          fullPath: tagPath.replace('//', '/'),
          affix: true
        })
      }
      if (route.children) {
        tags = tags.concat(filterAffixTags(route.children, route.path))
      }
    })
    return tags
  }
  const affixTags = filterAffixTags(routes)
  affixTags.forEach(tag => {
    if (tag.name) {
      tagsViewStore.addTagView(tag)
    }
  })
}

// 右键菜单
const menuVisible = ref(false)
const menuLeft = ref(0)
const menuTop = ref(0)
const selectedTag = ref<TagView | null>(null)

function openMenu(tag: TagView, e: MouseEvent) {
  const menuMinWidth = 105
  const offsetLeft = e.clientX + 15
  const offsetTop = e.clientY

  menuLeft.value = offsetLeft
  menuTop.value = offsetTop
  selectedTag.value = tag
  menuVisible.value = true
}

function closeMenu() {
  menuVisible.value = false
}

// 监听路由变化
watch(
  () => route.path,
  () => {
    addTags()
  }
)

// 点击其他地方关闭菜单
watch(menuVisible, (value) => {
  if (value) {
    document.body.addEventListener('click', closeMenu)
  } else {
    document.body.removeEventListener('click', closeMenu)
  }
})

onMounted(() => {
  initAffixTags()
  addTags()
})
</script>

<style lang="scss" scoped>
.tags-view-container {
  height: $tagsview-height;
  background-color: $bg-white;
  border-bottom: 1px solid $border-color-light;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.04);
  position: relative;
}

.tags-view-wrapper {
  :deep(.el-scrollbar__view) {
    display: flex;
    align-items: center;
    height: 100%;
    padding: 0 $spacing-sm;
    white-space: nowrap;
  }
}

.tags-view-item {
  display: inline-flex;
  align-items: center;
  height: 26px;
  padding: 0 $spacing-sm;
  margin: 0 2px;
  font-size: $font-size-small;
  color: $text-regular;
  background-color: $bg-white;
  border: 1px solid $border-color;
  border-radius: $border-radius-sm;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    color: $primary-color;
    border-color: $primary-color;
  }

  &.active {
    color: #fff;
    background-color: $primary-color;
    border-color: $primary-color;

    .close-icon {
      color: rgba(255, 255, 255, 0.8);

      &:hover {
        color: #fff;
        background-color: rgba(255, 255, 255, 0.2);
      }
    }
  }

  .close-icon {
    width: 16px;
    height: 16px;
    margin-left: 4px;
    font-size: 12px;
    border-radius: 50%;
    transition: all 0.2s;

    &:hover {
      background-color: rgba(0, 0, 0, 0.1);
    }
  }
}

.context-menu {
  position: fixed;
  z-index: 3000;
  margin: 0;
  padding: 5px 0;
  font-size: $font-size-small;
  color: $text-primary;
  list-style: none;
  background: $bg-white;
  border-radius: $border-radius-md;
  box-shadow: $shadow-light;

  li {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    padding: 8px $spacing-md;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
      background-color: $bg-hover;
    }
  }
}
</style>
