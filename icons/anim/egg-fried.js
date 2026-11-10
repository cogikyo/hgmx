import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const white = el.querySelector("path");
  animate(white, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const yolk = el.querySelector("circle");
  animate(yolk, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
