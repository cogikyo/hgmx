import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const scope = el.querySelectorAll("circle, path");
  animate(scope, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
  const ring = el.querySelector("circle");
  animate(ring, { scale: [1, 0.85, 1] }, { duration: τmid, delay: τmin });
}
