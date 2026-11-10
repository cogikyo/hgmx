import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const scope = el.querySelector("g");
  animate(scope, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
  const ticks = el.querySelector("g g");
  animate(ticks, { scale: [1, 1.35, 0.9, 1] }, { duration: τmid, delay: τmin });
}
