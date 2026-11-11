import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const reader = el.querySelector("g");
  animate(reader, { translateY: [0, 3, -1, 0] }, { duration: τmax });

  const hole = el.querySelector("circle");
  animate(hole, { scale: [1, 1.75, 1] }, { duration: τmin, delay: τmid });
}
