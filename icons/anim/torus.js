import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const outline = el.querySelector('ellipse[rx="10"]');
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hole = el.querySelector('ellipse[rx="3"]');
  animate(hole, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
