import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const strokes = el.querySelectorAll("path");
  animate(strokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const tail = el.querySelector("circle");
  animate(tail, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
