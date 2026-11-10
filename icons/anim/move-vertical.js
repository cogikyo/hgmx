import { animate, τmin } from "../motion.js";

export default function (el) {
  const nudge = { duration: 0.3 };
  animate(el.querySelector('path[d^="m8 6"]'), { translateY: [0, -3, 1, 0] }, { ...nudge, delay: 0 });
  animate(el.querySelector('path[d^="m8 18"]'), { translateY: [0, 3, -1, 0] }, { ...nudge, delay: τmin });
}
