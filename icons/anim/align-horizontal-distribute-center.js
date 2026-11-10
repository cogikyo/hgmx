import { animate, τmax } from "../motion.js";

export default function (el) {
  const before = el.querySelector('rect[x="4"]');
  animate(before, { translateX: [0, -1.5, 0.5, 0] }, { duration: τmax });

  const after = el.querySelector('rect[x="14"]');
  animate(after, { translateX: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
