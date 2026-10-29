import { animate } from "../motion.js";

export default function (el) {
  const nudge = { duration: 0.3 };
  animate(el.querySelector('path[d^="m9 5"]'), { translateY: [0, -3, 1, 0] }, { ...nudge, delay: 0 });
  animate(el.querySelector('path[d^="m19 9"]'), { translateX: [0, 3, -1, 0] }, { ...nudge, delay: 0.06 });
  animate(el.querySelector('path[d^="m15 19"]'), { translateY: [0, 3, -1, 0] }, { ...nudge, delay: 0.12 });
  animate(el.querySelector('path[d^="m5 9"]'), { translateX: [0, -3, 1, 0] }, { ...nudge, delay: 0.18 });
}
