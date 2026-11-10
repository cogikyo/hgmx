import { animate, τmax } from "../motion.js";

export default function (el) {
  const before = el.querySelector('rect[y="4"]');
  animate(before, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });

  const after = el.querySelector('rect[y="14"]');
  animate(after, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
