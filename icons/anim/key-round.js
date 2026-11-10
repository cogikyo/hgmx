import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const key = el.querySelectorAll("path, circle");
  animate(key, { translateX: [4, -1, 0], translateY: [-4, 1, 0] }, { duration: τmax });

  const hole = el.querySelector("circle");
  animate(hole, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
