export default defineAppConfig({
  ui: {
    primary: 'emerald',
    gray: 'slate',
    notifications: {
      position: 'top-right',
    },
    modal: {
      overlayTransition: {
        enter: 'ease-out duration-300',
        leave: 'ease-in duration-200',
      },
    },
    button: {
      default: {
        loadingIcon: 'i-heroicons-arrow-path',
      },
    },
    input: {
      default: {
        size: 'md',
      },
    },
    select: {
      default: {
        size: 'md',
      },
    },
    table: {
      default: {
        sortButton: {
          icon: 'i-heroicons-arrows-up-down',
        },
      },
    },
    pagination: {
      default: {
        activeButton: {
          color: 'primary',
        },
      },
    },
  },
})
