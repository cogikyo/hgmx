import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M21 7"]');
  const bottom = el.querySelector('path[d^="M21 17"]');
  animate(top, { translateY: [0, 2, -0.5, 0] }, { duration: τmid });
  animate(bottom, { translateY: [0, -2, 0.5, 0] }, { duration: τmid });

  const lid = el.querySelector('path[d^="M18.944"]');
  animate(lid, { scaleY: [1, 0.1, 1] }, { duration: τmid, delay: τmin });

  const pupil = el.querySelector("circle");
  animate(pupil, { scale: [1, 0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
