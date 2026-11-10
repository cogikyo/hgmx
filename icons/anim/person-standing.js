import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll("path");
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
