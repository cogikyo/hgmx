import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  animate(el.querySelector('rect[y="3"]'), { translateY: [0, -1.5, 0] }, { duration: τmax });
  animate(el.querySelector('rect[y="14"]'), { translateY: [0, 1.5, 0] }, { duration: τmax });

  const chevron = el.querySelector("path");
  animate(chevron, { translateX: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
