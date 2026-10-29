import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("path");
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const x = el.querySelectorAll("line");
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
