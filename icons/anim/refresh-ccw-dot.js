import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrows = el.querySelector("g");
  animate(arrows, { rotate: [0, -360], scale: [1, 0.85, 1] }, { duration: τmax });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
