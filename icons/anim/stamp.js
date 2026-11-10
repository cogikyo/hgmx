import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const stamp = el.querySelector("g");
  animate(stamp, { translateY: [0, 2.5, 0] }, { duration: τmax, times: [0, 0.35, 1] });

  const print = el.querySelector('path[d^="M5 22"]');
  animate(print, { scaleX: [0, 1] }, { duration: τmin, delay: τmin });
}
