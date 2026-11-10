import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  animate(el.querySelector('rect[x="3"]'), { translateX: [0, -1.5, 0] }, { duration: τmax });
  animate(el.querySelector('rect[x="14"]'), { translateX: [0, 1.5, 0] }, { duration: τmax });

  const chevron = el.querySelector("path");
  animate(chevron, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
