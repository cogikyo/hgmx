import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g.front");
  animate(front, { translateY: [0, -3, 1, 0], rotate: [0, -15, 6, 0] }, { duration: τmax });
  const back = el.querySelector("g.back");
  animate(back, { translateY: [0, -2, 0.5, 0], rotate: [0, 12, -4, 0] }, { duration: τmid, delay: τmin });
  const pop = { scale: [0, 1.5, 1], opacity: [0, 1] };
  animate(front.querySelectorAll('path[d$="h.01"]'), pop, { duration: τmid, delay: τmin });
  animate(back.querySelectorAll('path[d$="h.01"]'), pop, { duration: τmin, delay: τmid });
}
