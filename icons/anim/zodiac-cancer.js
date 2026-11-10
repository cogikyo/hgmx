import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arcs = el.querySelectorAll("path");
  animate(arcs, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const claws = el.querySelectorAll("circle");
  animate(claws, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
