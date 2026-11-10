import { animate, τmin } from "../motion.js";

export default function (el) {
  const nudge = { duration: 0.3 };
  animate(el.querySelector('path[d^="M19 13"]'), { translateX: [0, 2, -0.7, 0], translateY: [0, 2, -0.7, 0] }, { ...nudge, delay: 0 });
  animate(el.querySelector('path[d^="M5 11"]'), { translateX: [0, -2, 0.7, 0], translateY: [0, -2, 0.7, 0] }, { ...nudge, delay: τmin });
}
