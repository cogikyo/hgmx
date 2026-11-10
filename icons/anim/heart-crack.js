import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const heart = el.querySelector("g");
  animate(heart, { scale: [1, 1.25, 0.85, 1] }, { duration: τmid });
  const crack = el.querySelector('path[d^="M12.409"]');
  animate(crack, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
