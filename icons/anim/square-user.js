import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const body = el.querySelector("path");
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0], opacity: [1] }, { duration: τmin, delay: τmid });
}
