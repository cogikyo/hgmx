import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const disc = el.querySelector("circle");
  animate(disc, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 630] }, { duration: τmid });

  const hole = el.querySelector("path");
  animate(hole, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
