import Sortable from "sortablejs";

// 定义回调函数的类型，接收排序后的新旧索引
type SortCallback = (newIndex: number, oldIndex: number) => Promise<void>;

/**
 * 表格排序功能封装类
 * - 表格需要设置 row-key="id"
 * - 拖拽元素的类名只能是 drag-handler
 */
export class TableSortable {
  private instance: Sortable | null = null;

  /**
   * 表格排序功能封装类
   * - 表格需要设置 row-key="id"
   * - 拖拽元素的类名只能是 drag-handler
   * @param elSelector DOM 选择器 (例如 .el-table__body-wrapper tbody)
   * @param onSortEnd 排序完成后的回调函数
   */
  constructor(
    private elSelector: string,
    private onSortEnd: SortCallback,
  ) {}

  /** 初始化拖拽 */
  public init() {
    const el = document.querySelector(this.elSelector);

    if (!el) {
      console.warn(`TableSortable: 未找到元素 ${this.elSelector}`);
      return;
    }

    // 销毁旧实例防止重复绑定
    this.destroy();

    this.instance = Sortable.create(el as HTMLElement, {
      animation: 150,
      handle: ".drag-handler",
      onEnd: async ({ newIndex, oldIndex }) => {
        if (
          newIndex === oldIndex ||
          newIndex === undefined ||
          oldIndex === undefined
        )
          return;

        // 执行传入的业务逻辑
        try {
          await this.onSortEnd(newIndex, oldIndex);
        } catch (error) {
          console.error("拖拽排序回调执行失败:", error);
        }
      },
    });
  }

  /** 销毁实例 */
  public destroy() {
    if (this.instance) {
      this.instance.destroy();
      this.instance = null;
    }
  }
}
