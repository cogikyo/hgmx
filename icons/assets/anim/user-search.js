import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const searchGroup = el.querySelector("g.search");
  animate(searchGroup, { scale: [1.35, 1] }, { duration: τmax });
  animate(
    searchGroup,
    { translateX: [0, -12, 2, 0], translateY: [0, -6, -5, 0] },
    { duration: τmax },
  );

  const head = el.querySelector('circle[cx="10"]');
  animate(head, { translateY: [1, 1, -2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const body = el.querySelector('path[d*="15H7"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
