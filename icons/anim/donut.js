import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const donut = el.querySelector("path");
  animate(donut, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hole = el.querySelector("circle");
  animate(hole, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
