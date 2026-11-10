import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const palm = el.querySelector("g");
  animate(palm, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [0.5, 1.25, 1] }, { duration: τmid, delay: τmin });
}
