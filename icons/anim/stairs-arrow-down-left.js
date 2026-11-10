import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const steps = el.querySelector('path[d^="M2 20h5"]');
  animate(steps, { translateY: [0, 1, 0] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [0, -2, 1, 0], translateY: [0, 2, -1, 0] }, { duration: τmid, delay: τmin });
}
