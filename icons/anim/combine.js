import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [0, 2, -0.5, 0] }, { duration: τmid });

  const target = el.querySelector('rect[x="14"]');
  animate(target, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
