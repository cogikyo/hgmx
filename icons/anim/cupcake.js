import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wrapper = el.querySelectorAll('path[d^="M21 13"], path[d^="M12 22"], path[d^="m15.5"], path[d^="M8.5 22"]');
  animate(wrapper, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cherry = el.querySelector("circle");
  animate(cherry, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
